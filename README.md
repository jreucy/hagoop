HaGoop
======

A mapreduce framework built in Go with an emphasis on high churn


#### TCP Messages

J - Join Request (Worker to Server)  
M r [file\_name] [starting\_line] [ending\_line] - Map Request (Request to Server, Server to Worker)  
R r [file\_name] [starting\_line] [ending\_line] - Reduce Request (Request to Server, Server to Worker)  
M a [file\_name] [starting\_line] [ending\_line] - Map Answer (Worker to Server)  
R a [file\_name] [starting\_line] [ending\_line] - Reduce Answer (Worker to Server, Server to Request)  
X - Connection Lost (Worker to Server, Server to Request?)