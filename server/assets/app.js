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
        total: 0,
				filter: {
						"entry_type": 6,
						//"text": "الأندلس",
						//"tags": ["حِكمَة", "شِعر", "مَدح"],
						"limit": 50,
						"offset": 0
				}
    },
    mounted() {
        axios.post("/api/entry/query", this.filter)
            .then(response => {
                //console.log(response.data);
                this.total = response.data.total
                this.results = response.data.entries
            })
            .catch(function(error) {
                console.log(error);
            })
    },
    methods: {
        next: function() {
            if (this.filter.offset + this.filter.limit < this.total) {
                this.filter.offset = this.filter.offset + this.filter.limit
                axios.post("/api/entry/query", this.filter)
                    .then(response => {
                        //console.log(response.data);
                        this.total = response.data.total
                        this.results = response.data.entries
                    })
                    .catch(function(error) {
                        console.log(error);
                    })
            }
        },
        previous: function() {
            if (this.filter.offset - this.filter.limit >= 0) {
                this.filter.offset = this.filter.offset - this.filter.limit
                axios.post("/api/entry/query", this.filter)
                    .then(response => {
                        //console.log(response.data);
                        this.total = response.data.total
                        this.results = response.data.entries
                    })
                    .catch(function(error) {
                        console.log(error);
                    })
            }

        }
    }
});
