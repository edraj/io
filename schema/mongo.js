
db.content.createIndex( { displayname: "text", body: "text" }, {"default_language":"none"} ); // Sadly, Arabic language support requires 3rd party license
db.content.createIndex({files: 1});
db.content.createIndex({tags: 1});
db.file.createIndex({checksum: 1});
db.file.createIndex({tags: 1});
db.file.createIndex({shortname:"text", description:"text", text:"text"}) ;
db.file.createIndex({pathname:1});
db.message.createIndex({files:1});
