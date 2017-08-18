# Bird's-eye view
**Space center** manages **satellites**. Satellite does one thing and does it well.
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
