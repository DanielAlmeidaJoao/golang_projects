## each protocol has at least two go chans:
    * Message chan: this will only have the messages received
    * error chan: this will only have the errors sent to the protocol, or generic errors

## each protocol will have a map of functions that will handle the messages from those chans

## a server struct will have the following fields:
    * incomming message handler: this will process all the messages received
    * errorHandlerFunction: this will process all the received messages