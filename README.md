# Gogarin
For **software engineers** and **business automation engineers**, whose **apps** must be **integrated with 3rd-party asap**, the **Gogarin** is an **intelligent automation** tool that **integrates databases, services, apps** in a matter of minutes.


# Use Cases

### Tail -f

```
tail -f production.log
|-- filter message contains "sign-up"
|  |-- mailchimp subscribe {{user.id}} to list
|  |-- telegram send "Hi" to {{user.phone}}
|  |-- google_sheets append {{user.id}},{{user.email}},{{user.phone}} to "Users"
|-- postgresql insert {{user}} into table clients
```

### Issue Trackers

```
When a new issue is created
|-- filter {{issue.type}}="User Story"
|   |-- Create card in Trello using {{issue}}
|   |-- filter {{issue.assignee}}="Anton Kuzmenko" AND {{issue.priority}}="Critical"
|       |-- send "{{issue.key}} - {{issue.subject}}" to anton@email.com
```
