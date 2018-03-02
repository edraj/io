var messages = require('./edraj_pb');
var services = require('./edraj_grpc_pb');

var grpc = require('grpc');
const fs = require('fs');

function main() {

	const ssl_creds = grpc.credentials.createSsl(
    fs.readFileSync('../../../workspace/certs/edrajRootCA.crt'),
    fs.readFileSync('../../../workspace/certs/kefah.key'),
    fs.readFileSync('../../../workspace/certs/kefah.crt')
	);

	//var ssl_creds = grpc.credentials.createSsl(root_certs);
  var client = new services.OwnerClient('localhost:50051', ssl_creds); //grpc.credentials.createInsecure());
	console.log("Hi there");
	var filter = new messages.Filter();
	filter.setEntryType(messages.EntryType.CONTENT);
	filter.setLimit(5);
	client.query(filter, function(err, res) {
	console.log('Res:', res.getTotal());
	entries = res.getEntriesList();	
	for(let i=0; i< entries.length; i++) {
		one = entries[i].getContent();
		console.log("Entry: ", i, one.getDisplayname())
	}
	});
}

main();
