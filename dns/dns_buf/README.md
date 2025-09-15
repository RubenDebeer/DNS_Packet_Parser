## Parsing the DNS Packet


#### Headers

**![Header](./dns/images/Header_Format.png)**

| RFC Name | Descriptive Name | Length | Description |
| ----- | ----- | ----- | ----- |
| ID | Packet Identifier | 16 bits ( 2 Bytes) | A random identifier is assigned to query packets. Response packets must reply with the same id. This is needed to differentiate responses due to the stateless nature of UDP. |
| QR | Query Response | 1 bit | 0 for queries, 1 for responses. |
| OPCODE | Operation Code | 4 bits | **Typically always 0**, see RFC1035 for details. |
| Flag :AA | Authoritative Answer | 1 bit | Set to 1 if the responding server is authoritative \- that is, it "owns" \- the domain queried. |
| Flag :TC | Truncated Message | 1 bit | Set to 1 if the message length exceeds 512 bytes. Traditionally a hint that the query can be reissued using TCP, for which the length limitation doesn't apply. |
| Flag :RD | Recursion Desired | 1 bit | Set by the sender of the request if the server should attempt to resolve the query recursively if it does not have an answer readily available. |
| Flag :RA | Recursion Available | 1 bit | Set by the server to indicate whether or not recursive queries are allowed. |
| Flag :Z | Reserved | 3 bits | Originally reserved for later use, but now used for DNSSEC queries. |
| Flag :RCODE | Response Code | 4 bits | Set by the server to indicate the status of the response, i.e. whether or not it was successful or failed, and in the latter case providing details about the cause of the failure. |
| QDCOUNT | Question Count | 16 bits | The number of entries in the Question Section |
| ANCOUNT | Answer Count | 16 bits | The number of entries in the Answer Section |
| NSCOUNT | Authority Count | 16 bits | The number of entries in the Authority Section |
| ARCOUNT | Additional Count | 16 bits | The number of entries in the Additional Section |


#### Questions
**![Questions](./dns/images/Question_Format.png)**

Standard Query
A Standard query specifies a target domain name (QNAME) , query type (QTYPE) and a query class (QCLASS) and asks for the matching RR. The QTYPE and QCLASS fields are each 16 bits long (2bytes).

| Field | Type | Description |
| ----- | ----- | ----- |
| Name | Label Sequence | The **domain name**, encoded as a sequence of labels.(QNAME) |
| Type | 2-byte Integer | The record type. |
| Class | 2-byte Integer | The class, in practice, is always set to 1\. |

#### Answers
**![Answers](./dns/images/Answer_Format.png)**

| Field | Type | Description |
| ----- | ----- | ----- |
| Name | Label Sequence | The domain name, encoded as a sequence of labels. |
| Type | 2-byte Integer | The record type. |
| Class | 2-byte Integer | The class, in practice, is always set to 1\. |
| TTL | 4-byte Integer | Time-To-Live, i.e. how long a record can be cached before it should be required. |
| Len | 2-byte Integer | Length of the record type specific data. |