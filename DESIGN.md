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

## Tail

```
tail -f production.log
|-- filter message contains "sign-up"
|  |-- mailchimp subscribe {{user.id}} to list
|  |-- telegram send "Hi" to {{user.phone}}
|  |-- google_sheets append {{user.id}},{{user.email}},{{user.phone}} to "Users"
|-- postgresql insert {{user}} into table clients
```

## JIRA

```
When new issues is created
|-- filter {{issue.type}}="User Story"
|   |-- Create card in Trello using {{jira_issue}}
|   |-- filter {{issue.assignee}}=Anton Kuzmenko AND {{issue.priority}}=Critical
|       |-- send "{{issue.key}} - {{issue.subject}}" to anton@email.com
```