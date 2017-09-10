[![Go Report Card](https://goreportcard.com/badge/github.com/antonkuzmenko/gogarin)](https://goreportcard.com/report/github.com/antonkuzmenko/gogarin)
[![Build Status](https://travis-ci.org/antonkuzmenko/gogarin.svg?branch=master)](https://travis-ci.org/antonkuzmenko/gogarin)
[![MIT license](https://img.shields.io/badge/license-MIT%202.0-blue.svg?style=flat-square)](https://github.com/antonkuzmenko/gogarin/blob/master/LICENSE)

# Gogarin
For software engineers and business automation engineers who want to integrate various infrastructure components, Gogarin is an intelligent automation tool that allows you to integrate databases, services, and applications in a matter of minutes.


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

# It Is

1. Fast
1. Scalable
1. Persistent
1. Fault-tolerant
1. Simple
1. Lightweight
1. Easy to install
1. Easy to maintain
