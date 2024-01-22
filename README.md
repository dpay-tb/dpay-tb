# dpay-tb
This is a rewrite of DPay using TigerbeetleDB. Implementing financial transactions in a traditional relational database may pose challenges such as ensuring ACID compliance under high concurrency, addressing potential concurrency issues, avoiding performance bottlenecks with careful indexing and query optimization, and maintaining data integrity through thorough validation processes. Turns out there's a DB which builds financial transactions into it as a domain specific data structure: [TigerbeetleDB](https://tigerbeetle.com/).

This project intends to use this DB to handle financial transactions, and hopefully reap performance benefits and improved correctness and data integrity. Right now it is very rudimentary, and we hope you may be able to contribute more features to it!

# Setting up Tigerbeetle
We provide a simple Dockerfile in the `tigerbeetle/` directory. However we recommend you to download and run a local version. You may follow the [official instructions](https://docs.tigerbeetle.com/quick-start/single-binary). The project uses the Go client, so check out the [docs](https://docs.tigerbeetle.com/clients/go).
