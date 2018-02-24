#!/bin/bash
mongo edraj --eval 'db.dropDatabase();'  && ./maqola_loader.py && mongo edraj ../../schema/mongo.js
