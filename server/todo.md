Implement:

 1. Attachment creation and retrieval /file upload/downaload [chunked]: for a container we can use tar.gz.
 2. Update including attachment
 3. General event notification (should be stored to mongo: notifications)
 4. Messaging (full test cycle) with notification
 5. Android mobile app 
 6. Layout, Page and Block: SPA + SEO Friendly
 7. Addon
 8. Schema example + structured content body
 9. Federation (server2server)
10. Multi-device per user (multi signed certificates for the same commonName)
11. One mongo database per user?
12. Bit-torrent-like large file serving (>10mb)
13. Write arabic language support for mongodb
14. Buffer/cache undelivered requests until the respective domain server is up again (2-4 days max?)

15. Consider https://github.com/cloudflare/cfssl instead of certstrap. cfss has: 5. revoke, 6. intermediate ca
certstrap can 1. create CA, 2. Create requests, 3. distinguish between domain (server) and user (client), 4. Sign. 
```bash
# main tool
go get -u github.com/cloudflare/cfssl/cmd/cfssl
# Other tools
go get -u github.com/cloudflare/cfssl/cmd/...
```

16. Should consul be additionally used by server (in addition to root?) in that case can we enable strict two-way ssl verfication?

https://www.programming-books.io/essential/go/
