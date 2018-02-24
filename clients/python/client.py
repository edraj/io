#!/bin/env python3
"""The Python implementation of the GRPC edraj.EntryService client."""

from __future__ import print_function

import grpc

import edraj_pb2
import edraj_pb2_grpc

def run():

    with open('../../../workspace/certs/edrajRootCA.crt') as f:
        rootca_crt = bytes(f.read(), 'utf8')

    with open('../../../workspace/certs/admin.key') as f:
        kefah_key = bytes(f.read(), 'utf8')

    with open('../../../workspace/certs/admin.crt') as f:
        kefah_crt = bytes(f.read(), 'utf8')


    credentials = grpc.ssl_channel_credentials(root_certificates=rootca_crt, private_key=kefah_key, certificate_chain=kefah_crt )
    channel = grpc.secure_channel('localhost:50051', credentials)
    stub = edraj_pb2_grpc.OwnerStub(channel)
    response = stub.Query(edraj_pb2.Filter(entry_type='CONTENT'), metadata=[('edraj-foo', 'bar')])

    for entry in response.entries:
        #print(entry.content.tags)
        print(entry.content)


if __name__ == '__main__':
    run()
