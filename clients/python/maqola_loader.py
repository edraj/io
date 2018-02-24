#!/bin/env python3
"""The Python implementation of the GRPC edraj.EntryService client."""

from __future__ import print_function

import grpc

import edraj_pb2 as edraj
import edraj_pb2_grpc
import pymysql.cursors
from datetime import datetime
import subprocess
import os
import sys

def load(stub):
    mysql = pymysql.connect(host='localhost', user='maqola', password='maqola', db='maqola', charset='utf8mb4', 
                            cursorclass=pymysql.cursors.DictCursor)

    cursor = mysql.cursor()
    sub_cursor = mysql.cursor()
    users = {}
    authors = {}
    fileids = []
    cursor.execute("select * from tbl_user")
    all = cursor.fetchall()
    for one in all: 
        users[one['id']] = {'prettyname': one['prettyname']}
    
    cursor.execute("select * from tbl_author")
    all = cursor.fetchall()
    for one in all: 
        authors[one['id']] = {'prettyname': one['prettyname'], 'description': one['description'], 'shortname': one['shortname']}
    
    cursor.execute("select * from tbl_entry")
    all = cursor.fetchall()
    for one in all:
        entry = edraj.Content ( 
            id = str(one['id']), 
            displayname = one['title'], 
            created = one['created_at'], 
            body = one['text'], 
            author = edraj.Identity(id=str(one['author_id'])),
            actor = edraj.Identity(id=str(one['user_id'])))
        if users[one['user_id']]:
            entry.actor.displayname = users[one['user_id']]['prettyname']

        if authors[one['author_id']]:
            entry.author.displayname = authors[one['author_id']]['prettyname']
            entry.author.shortname   = authors[one['author_id']]['shortname']
            entry.author.description = authors[one['author_id']]['description']
        
        if one['updated_at']:
            entry.updated = one['updated_at']
        
        sub_cursor.execute("select * from tbl_entry_media where entry_id = %s", (one['id']))
        recs = sub_cursor.fetchall()
        for rec in recs:
            # attachment = {'file_uri': rec['filename'], 'status': rec['publish_status']}
            # attachment = entry.attachments.add(file_path = str(one['id']%100) + '/' + str(one['id']) + '/' + rec['filename'])
            entry.files.add(str(rec['id']))
            if rec['id'] in fileids:
                print("skipping file id %d" % rec['id'])
                break
            fileids.append(rec['id'])
            file = edraj.File(id=str(rec['id']), pathname = str(one['id']%100) + '/' + str(one['id']) + '/' + rec['filename'])
            if rec['description']:
                file.description = rec['description']
            file.shortname = rec['filename']
            file.type = edraj.File.JPG if rec['type'] == 1 else edraj.File.MP3 
            full_pathname = '../../../workspace/media/' + file.pathname
            file.mime = subprocess.run(['file', '-b', '-i', full_pathname], stdout=subprocess.PIPE).stdout.decode('utf-8').strip()
            file.mime_description = subprocess.run(['file', '-b', full_pathname], stdout=subprocess.PIPE).stdout.decode('utf-8').strip()
            file.size =  os.path.getsize(full_pathname)
            file.checksum = "sha256:" + subprocess.run(['sha256sum', '-b', full_pathname], stdout=subprocess.PIPE).stdout.decode('utf-8').strip().partition(' ')[0]
            file.storage = edraj.File.PATHNAME
            response = stub.Create(edraj.Entry(file = file, id = entry.id, type = edraj.FILE))

        sub_cursor.execute("select * from tbl_entry_tag where entry_id = %s", (one['id']))
        recs = sub_cursor.fetchall()
        for rec in recs:
            entry.tags.append(rec['tag_id'])
        
        sub_cursor.execute("select * from tbl_entry_comment where entry_id = %s", (one['id']))
        recs = sub_cursor.fetchall()
        for rec in recs:                     
            comment = entry.comments.add(
                actor=edraj.Identity(id=str(rec['user_id'])), 
                created= rec['created_at'], 
                #'status': rec['publish_status'], 
                body= rec['text'])
            if users[rec['user_id']] and users[rec['user_id']]['prettyname']:
                comment.actor.displayname = users[rec['user_id']]['prettyname']
            if rec['title']:
                comment.title = rec['title']
            if rec['updated_at']:
                comment.updated = rec['updated_at']
        
        sub_cursor.execute("select * from tbl_entry_vote where entry_id = %s", (one['id']))
        recs = sub_cursor.fetchall()
        for rec in recs:
            reaction = entry.reactions.add(actor=edraj.Identity(id=str(rec['user_id'])))
            if rec['created_at']:
                reaction.created = rec['created_at']
            if rec['updated_at']:
                reaction.updated = rec['updated_at']
            reaction.type = edraj.Reaction.LIKE if rec['type'] == 1 else edraj.Reaction.DISLIKE
        

        """
        sub_entry = []
        sub_cursor.execute("select * from tbl_entry_tag where entry_id = %s", (one['id']))
        recs = sub_cursor.fetchall()
        for rec in recs:
            sub_entry.append({})
        if sub_cursor:
            entry['views'] = sub_cursor
        

        changes = []
        sub_cursor.execute("select * from tbl_entry_update_log where entry_id = %s", (one['id']))
        recs = sub_cursor.fetchall()
        for rec in recs:
            change = {
                "actor": rec['user_id'],
                "old_id": rec['id'],
                "old": {}
                }

            if rec['updated_at']:
                change["created_at"] = datetime.fromtimestamp(rec['updated_at'])

            if rec['old_title']:
                change['old']['title'] = rec['old_title']

            if rec['old_author_id']:
                change['old']['actor'] = rec['old_author_id']

            if rec['old_tags']:
                change['old']['tags'] = rec['old_tags']

            if rec['old_medias']:
                change['old']['media'] = rec['old_medias']

            if not change['old']:
                del change['old']

            changes.append(change)

        if changes:
            entry['history'] = changes

        """
        response = stub.Create(edraj.Entry(content = entry, id = entry.id, type = edraj.CONTENT))
        print(entry.id, response)
        # print("Inserted %r : %r" % (one['id'], entry['_id']))


def run():
    workspace = '../../../workspace'
    with open(workspace + '/certs/edrajRootCA.crt') as f:
        rootca_crt = bytes(f.read(), 'utf8')

    with open(workspace + '/certs/admin.key') as f:
        user_key = bytes(f.read(), 'utf8')

    with open(workspace + '/certs/admin.crt') as f:
        user_crt = bytes(f.read(), 'utf8')

    credentials = grpc.ssl_channel_credentials(root_certificates=rootca_crt, private_key=user_key, certificate_chain=user_crt )

    channel = grpc.secure_channel('localhost:50051', credentials)
    stub = edraj_pb2_grpc.OwnerStub(channel)

    load(stub)


if __name__ == '__main__':
    run()
