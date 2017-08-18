# Vision

For **software engineers** and **business automation engineers**, whose **apps** must be **integrated with 3rd-party asap**, the **Gogarin** is an **intelligent automation** tool that **integrates databases, services, apps, and so forth** in a matter of minutes.

# It Is

1. Fast
1. Scalable
1. Persistent
1. Fault-tolerant
1. Simple
1. Lightweight
1. Easy to install
1. Easy to maintain

# Use Cases

### Tail

```
tail -f production.log
|-- filter message contains "sign-up"
|  |-- mailchimp subscribe {{user.id}} to list
|  |-- telegram send "Hi" to {{user.phone}}
|  |-- google_sheets append {{user.id}},{{user.email}},{{user.phone}} to "Users"
|-- postgresql insert {{user}} into table clients
```

### JIRA

```
When new issue is created
|-- filter {{issue.type}}="User Story"
|   |-- Create card in Trello using {{jira_issue}}
|   |-- filter {{issue.assignee}}="Anton Kuzmenko" AND {{issue.priority}}="Critical"
|       |-- send "{{issue.key}} - {{issue.subject}}" to anton@email.com
```

# Bird's-eye view
Space center manages satelites.
Satellite does one thing and does it well.
Example satelites:
 - Triggers
   - Tail trigger (`tail -f somefile`)
   - HTTP callback
   - Cron trigger
   - JIRA trigger
   - POP3 trigger (receive emails)
   - Twitter trigger (react on a #hashtag)
   - and more.
 - Filters
   - Basic filter (`field == val`, `number >= value`, etc.)
   - Postgresql filter (`SELECT`)
   - Redis filter (`GET`, `HEXISTS`, `HGET`, etc.)
   - etc.
 - Modifiers
   - Basic modifier (add/delete/change fields of a message)
   - Postgresql modifier (add/delete/change fields of a message based on a `SELECT`)
   - JIRA modifier (add/delete/change fields based on JIRA API query results)
   - etc.
 - Actions
   - File action (append message to a file, create/update/delete file)
   - Postgresql action (`INSERT`, `UPDATE`, `DELETE`)
   - Mail action (send email)
   - Twitter action (post a tweet)
   - JIRA action (create/update/transition/etc. an issue)
   - etc.
 - Splitters
   - Basic splitter (iterates over a collection and triggers each message)
