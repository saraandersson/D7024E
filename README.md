# D7024E - Kademlia
## Compile and run the code
- cd d7024e
- docker build . -t kadlab 
- docker-compose up --scale kademliaNodes=49
### To enter a node: 
- cd d7024e
- docker attach d7024e_kademliaNodes_1 <br>
If you would like to enter another node then just change the node number 1 to a number between 1-49.
### Cli-commands inside a kademlia node
- put
- Enter data to put: (This could for example be "test")
- get
- Enter key to get data: (Here the filekey should be entered)
- exit
