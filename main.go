package main

//DNS packet structure
//+---------------------+
//|        Header       |
//+---------------------+
//|       Question      | Question for the name server
//+---------------------+
//|        Answer       | records that were requested from above
//+---------------------+
//|      Authority      | Name servers that resolve queries recursively
//+---------------------+
//|      Additional     | additional records

// DNS HEADER
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+---------+
//|  ID(16bit) - assigned  randomly                      |
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+---------+
//|QR(1b)|Opcode(4b)|AA(4b)|TC|RD|RA|Z    |   RCODE      |
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+----------
//|                    QDCOUNT                           |
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+---------+
//|                    ANCOUNT                           |
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+---------+
//|                    NSCOUNT                           |
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+------+
//|                    ARCOUNT                           |
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+------+
