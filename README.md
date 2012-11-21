HaGoop
======

A mapreduce framework built in Go with an emphasis on high churn


#### TCP Messages

J - Join Request (Worker to Server) 

MR [file\_name] [starting\_line] [ending\_line] - MapReduce Request (Request to Server)  
M [file\_name] [starting\_line] [ending\_line] - Map Request (Server to Worker)  
R [file\_name] [starting\_line] [ending\_line] - Reduce Request (Server to Worker)  
Ma [file\_name] [starting\_line] [ending\_line] - Map Answer (Worker to Server)  
Ra [file\_name] [starting\_line] [ending\_line] - Reduce Answer (Worker to Server, Server to Request)  
X - Connection Lost (Worker to Server, Server to Request?)