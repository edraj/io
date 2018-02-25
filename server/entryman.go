package main

import (
	"fmt"
	"path"
	"reflect"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	trash = "trash"
	edraj = "edraj"
)

// EntryMan ...
type EntryMan struct {
	mongoSession *mgo.Session
	mongoDb      *mgo.Database
	fileStore    Storage
}

func (man *EntryMan) init(config *Config) (err error) {
	man.mongoSession, err = mgo.Dial(config.mongoAddress)
	if err != nil {
		return
	}
	man.mongoDb = man.mongoSession.DB(edraj)
	man.fileStore.RootPath = path.Join(config.dataPath, edraj)
	man.fileStore.TrashPath = path.Join(config.dataPath, trash, edraj)
	return
}

func entryObject(objectType string, e *Entry, createIfNil bool) (doc interface{}) {
	fieldName := strings.Title(objectType)
	field := reflect.ValueOf(e).Elem().FieldByName(fieldName)
	doc = field.Interface()
	if createIfNil && field.IsNil() {
		field.Set(reflect.Indirect(reflect.New(field.Type().Elem())).Addr())
		doc = field.Interface()
	}
	return
}

func (man *EntryMan) create(entry *Entry) (response *Response, err error) {
	response = &Response{Status: &Status{}}
	//glog.Info(request)
	/*if entry.Type != nil {
		response.Status.Code = int32(codes.InvalidArgument)
		err = status.Errorf(codes.InvalidArgument, "Entry details are missing (%v).", request.Entry)
		response.Status.Message = err.Error()
		return
	}*/
	entryType := strings.ToLower(entry.Type.String())

	doc := entryObject(entryType, entry, false)

	err = man.mongoDb.C(entryType).Insert(doc)
	if err != nil {
		response.Status.Code = int32(codes.Internal)
		response.Status.Message = fmt.Sprintf("Failed to create entry (%v).", entry)
		err = status.Errorf(codes.Internal, err.Error())
		return
	}

	response.Status.Code = int32(codes.OK)
	return
}

func (man *EntryMan) update(entry *Entry) (response *Response, err error) {
	response = &Response{Status: &Status{}}
	//glog.Info(request)
	if /*entry.Type == "" ||*/ entry.Id == "" {
		response.Status.Code = int32(codes.InvalidArgument)
		err = status.Errorf(codes.Internal, "Entry details are missing (%v).", entry)
		response.Status.Message = err.Error()
		return
	}

	entryType := strings.ToLower(entry.Type.String())

	doc := entryObject(entryType, entry, false)

	err = man.mongoDb.C(entryType).Update(entry.Id, doc)
	if err != nil {
		response.Status.Code = int32(codes.NotFound)
		response.Status.Message = fmt.Sprintf("Failed to update entry (%v).", entry)
		err = status.Errorf(codes.NotFound, err.Error())
		return
	}

	response.Status.Code = int32(codes.OK)
	return
}

func (man *EntryMan) query(filter *Filter) (response *Response, err error) {
	response = &Response{Status: &Status{}}
	//glog.Info(request)

	/*if query.Type == "" {
		response.Status.Code = int32(codes.InvalidArgument)
		err = status.Errorf(codes.InvalidArgument, "Filter details are missing (%v).", query)
		response.Status.Message = err.Error()
		return
	}*/
	// TODO consume the Query data (filters, pagination ...etc)
	query := bson.M{}

	response.Entries = []*Entry{}
	entryType := strings.ToLower(filter.EntryType.String())

	if filter.Text != "" {
		query["$text"] = bson.M{"$search": filter.Text}
	}

	if len(filter.EntryIds) > 0 {
		query["_id"] = bson.M{"$in": filter.EntryIds}
	}

	if len(filter.Tags) > 0 {
		query["tags"] = bson.M{"$all": filter.Tags} // $all for and, $in for or
	}

	grpclog.Info("Filter: ", filter)
	grpclog.Info("Query: ", query)

	fieldName := strings.Title(entryType)
	fieldType := reflect.New(reflect.TypeOf(Entry{})).Elem().FieldByName(fieldName).Type().Elem()
	slice := reflect.MakeSlice(reflect.SliceOf(fieldType), 0, 0)
	objects := reflect.New(slice.Type())
	objects.Elem().Set(slice)
	if filter.Limit == 0 {
		filter.Limit = 2
	}
	err = man.mongoDb.C(entryType).Find(query).Skip(int(filter.Offset)).Limit(int(filter.Limit)).All(objects.Interface())
	if err != nil {
		response.Status.Code = int32(codes.Internal)
		err = status.Errorf(codes.Internal, err.Error())
		return
	}

	for i := 0; i < objects.Elem().Len(); i++ {
		entry := &Entry{}
		reflect.ValueOf(entry).Elem().FieldByName(fieldName).Set(reflect.ValueOf(objects.Elem().Index(i).Addr().Interface()))
		response.Entries = append(response.Entries, entry)
	}

	//response.Returned = int64(len(response.Entries))
	total, _ := man.mongoDb.C(entryType).Find(query).Count()
	response.Total = int64(total)

	response.Status.Code = int32(codes.OK)
	grpclog.Info("Response code: ", response.Status, "count: ", response.Total)
	return
}

/*
func (man *EntryMan) get(request *Entry) (response *Response, err error) {
	//glog.Info(request)
	response = &Response{Status: &Status{}}

	if request.EntryId == "" {
		response.Status.Code = int32(codes.InvalidArgument)
		err = status.Errorf(codes.InvalidArgument, "EntryId (%v) or EntryType (%v) are missing.", request.EntryId, request.EntryType)
		response.Status.Message = err.Error()
		return
	}

	entryType := strings.ToLower(request.EntryType.String())
	entry := Entry{}
	doc := entryObject(entryType, &entry, true)

	err = man.mongoDb.C(entryType).FindId(request.EntryId ).One(doc) //bson.ObjectIdHex(request.EntryID)
	if err != nil {
		response.Status.Code = int32(codes.NotFound)
		response.Status.Message = fmt.Sprintf("Failed to get entry (%v).", request.EntryId)
		err = status.Errorf(codes.NotFound, err.Error())
		return
	}

	response.Status.Code = int32(codes.OK)
	response.Entries = []*Entry{&entry}
	response.Returned = 1
	response.Total = 1

	return
}
*/

func (man *EntryMan) delete(entry *Entry) (response *Response, err error) {
	response = &Response{Status: &Status{}}
	if /*entry.Type == "" ||*/ entry.Id == "" {
		response.Status.Code = int32(codes.InvalidArgument)
		err = status.Errorf(codes.InvalidArgument, "EntryId (%v) or EntryType (%v) are missing.", entry.Id, entry.Type)
		response.Status.Message = err.Error()
		return
	}

	entryType := strings.ToLower(entry.Type.String())
	err = man.mongoDb.C(entryType).Remove(&struct{ _id string }{_id: entry.Id})
	if err != nil {
		response.Status.Code = int32(codes.NotFound)
		response.Status.Message = fmt.Sprintf("Failed to delete entry (%v).", entry.Id)
		err = status.Errorf(codes.NotFound, err.Error())
		return
	}

	response.Status.Code = int32(codes.OK)
	return
}
