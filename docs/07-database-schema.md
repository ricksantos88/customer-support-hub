# Database Schema

## contacts

```sql
CREATE TABLE contacts (
  id UUID PRIMARY KEY,
  phone VARCHAR(20) UNIQUE,
  name VARCHAR(255),
  created_at TIMESTAMP
);
```

## conversations

```sql
CREATE TABLE conversations (
  id UUID PRIMARY KEY,
  contact_id UUID REFERENCES contacts(id),
  status VARCHAR(20),
  assigned_agent_id UUID,
  last_message_at TIMESTAMP
);
```

## messages

```sql
CREATE TABLE messages (
  id UUID PRIMARY KEY,
  conversation_id UUID REFERENCES conversations(id),
  direction VARCHAR(20),
  content TEXT,
  external_id VARCHAR(255),
  status VARCHAR(20),
  created_at TIMESTAMP
);
```

Indexes:
- conversation_id
- phone
- created_at