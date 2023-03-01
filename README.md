# Chat App Testing

## Requirements
-   High ingress
-   Eventual consistency between participants
-   Consistency < latency

## Where is the data stored?

-   Db, textfiles, some datastore on the backend
-   It is not a transfer mechanism
-   Data stored in batches

## Application Layer

-   Manages communication between participants
-   Memory used as transfer mechanism
-   Retrieval mechanism for historic data (which could be paginated) or reconciliation
	-   This could be sent from memory or datastore
-   Transient within this layer
-   How much memory available can determine how much data is stored
-   LRU memory buffer
-   Node / Go - something that handles concurrency and streams well

## Presentation Layer

-   JS
-   Websockets

### Entities

-   Conversations
	-   Collection of messages
-   Messages
	-   Created by participant
	-   Has to belong to a conversation
-   Participants
	-   Has to belong to one or more conversations
	-   Can create messages within conversations
	-   Can join and leave conversations

## Participants

-   Unique identifier
	-   GUID - Possibly username - cellphone number
	-   For external use
-   Internal identifier
	-   Int as primary key
-   Name
	-   Displayed username

## Conversations

-   Internal identifier
-   Name
-   Creation date
-   Global flag for history visibility

## Conversation_Participants Join Table

-   Conversation_id
-   Participant_id
-   Join date

## Messages

-   Internal identifier (Message id)
-   Conversation_id
-   Participant_id
-   Message body
-   Timestamp

## Use cases (Lean inception)

-   Create Profile/Participant
	-   If not in memory, loads into database
	-   After creation, add to memory store
-   Create Conversations
	-   Find by username
	-   Starting channel
	-   Implicitly joins conversation
-   Join Conversation
	-   Registering to a channel
	-   Joining strea
