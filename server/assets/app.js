/*const vm = new Vue({
    el: '#app',
    data: {
        results: [
            { title: "the very first post", abstract: "lorem ipsum some test dimpsum" },
            { title: "and then there was the second", abstract: "lorem ipsum some test dimsum" },
            { title: "third time's a charm", abstract: "lorem ipsum some test dimsum" },
            { title: "four the last time", abstract: "lorem ipsum some test dimsum" }
        ]
    }
});*/



const vm = new Vue({
    el: '#app',
    data: {
        results: [],
        attachments: {},
        total: 0,
        filter: {
            "entry_type": 6,
            //"entry_ids": ["422"],
            //"text": "الأندلس",
            //"tags": ["حِكمَة", "شِعر", "مَدح"],
            "limit": 50,
            "offset": 0
        }
    },
    /*mounted()*/
    created() {
        this.load()
    },
    methods: {
        load: function() {
            axios.post("/api/entry/query", this.filter)
                .then(response => {
                    var l = response.data.entries.length;
                    fileids = new Set();
                    for (var i = 0; i < l; i++) {
                        if (response.data.entries[i].content.hasOwnProperty('files')) {
                            response.data.entries[i].content.files.forEach(e => {
                                fileids.add(e);
                            });
                        }
                    }

                    if (fileids.size > 0) {
                        axios.post("/api/entry/query", { "entry_type": 2, "offset": 0, "limit": 10, "entry_ids": [...fileids] })
                            .then(r => {
                                r.data.entries.forEach(e => {
                                    this.attachments[e.file.id] = e.file;
                                });
                            })
                            .catch(function(error) {
                                console.log(error);
                            })
                    }
                    //console.log(this.attachments);
                    //new Promise((resolve) => setTimeout(resolve, 500));
                    this.total = response.data.total;
                    this.results = response.data.entries;
                })
                .catch(function(error) {
                    console.log(error);
                })

        },
        next: function() {
            if (this.filter.offset + this.filter.limit < this.total) {
                this.filter.offset = this.filter.offset + this.filter.limit;
                this.load();
            }
        },
        previous: function() {
            if (this.filter.offset - this.filter.limit >= 0) {
                this.filter.offset = this.filter.offset - this.filter.limit
                this.load();
            }

        }
    }
});