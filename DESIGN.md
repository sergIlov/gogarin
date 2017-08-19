# Bird's-eye view
**Space center** manages **satellites**. Satellite does one thing and does it well.
A satellite is a Go application that communicates with the space center via Redis/RabbitMQ/NSQ/Apache Kafka.
Satellite is either a trigger, a filter, a modifier, an action, or a splitter.
You can maximize throughput by running many instances of the same satellite on different hosts.

```
              ++++++++++++++
             |    host1     |
             |              |
             | Space center |
             |              |
              ++++++++++++++
                     ▲
                     |
                     ▼
         |--------►NSQD◄-------|
         ▼                     ▼
 ++++++++++++++++      ++++++++++++++++
|     host2      |    |     host3      |
|                |    |                |
| Basic modifier |    | Basic modifier |
| Basic splitter |    | Basic splitter |
| Basic filter   |    | Basic filter   |
| Redis filter   |    | Redis action   |
| Tail trigger   |    | File action    |
| Cron trigger   |    | Cron trigger   |
| JIRA trigger   |    | JIRA modifier  |
| Mail action    |    | Mail action    |
 ++++++++++++++++      ++++++++++++++++
 ````

### Satelites
Example satellites:
 - **Triggers**
   - Tail trigger (`tail -f somefile`)
   - HTTP callback trigger
   - Cron trigger
   - JIRA trigger
   - POP3 trigger (receive emails)
   - Twitter trigger (react on a #hashtag)
   - and more.
 - **Filters**
   - Basic filter (`field == val`, `number >= value`, etc.)
   - PostgreSQL filter (`SELECT`)
   - Redis filter (`GET`, `HEXISTS`, `HGET`, etc.)
   - etc.
 - **Modifiers**
   - Basic modifier (add/delete/change fields of a message)
   - PostgreSQL modifier (add/delete/change fields of a message based on a `SELECT`)
   - JIRA modifier (add/delete/change fields based on JIRA API query results)
   - etc.
 - **Actions**
   - File action (append message to a file, create/update/delete file)
   - PostgreSQL action (`INSERT`, `UPDATE`, `DELETE`)
   - Mail action (send email)
   - Twitter action (post a tweet)
   - JIRA action (create/update/transition/etc. an issue)
   - etc.
 - **Splitters**
   - Basic splitter (iterates over a collection and triggers each item as a message)
