# Pento tech challenge

Thanks for taking the time to do our tech challenge. 

The challenge is to build a small full stack web app, that can help a freelancer track their time.

It should satisfy these user stories:

- As a user, I want to be able to start a time tracking session
- As a user, I want to be able to stop a time tracking session
- As a user, I want to be able to name my time tracking session
- As a user, I want to be able to save my time tracking session when I am done with it
- As a user, I want an overview of my sessions for the day, week and month
- As a user, I want to be able to close my browser and shut down my computer and still have my sessions visible to me when I power it up again.

## Getting started

You can fork this repo and use the fork as a basis for your project. We don't have any requirements on what stack you use to solve the task, so there is nothing set up beforehand.

## Timing

- Don't spend more than a days work on this challenge. We're not looking for perfection, rather try to show us something special and have reasons for your decisions.
- Get back to us when you have a timeline for when you are done.

## Notes

- This is technically possible to implement only on the frontend, but please take the opportunity to show your skills on the entire stack 
- Please focus more on code quality, building a robust service and such, than on the UI.

## Implementation details

Using Postgres as datastore
User will be presented by uuid

## Logic

Main domain logic will be described by datastore structure.

Using Postgres and single table.

Attempting to:

1. Avoid expensive JOIN
2. Possibility to migrate into NoSQL with ease if there would be future change in datastore engine
3. Multitenant system, using `user_id` i.e. eventually more than one user shoud be able to use the app

```sql
CREATE TABLE public.time_session_partial
(
    time_session_partial_id uuid NOT NULL,
    time_session_name character varying(128) COLLATE pg_catalog."default" NOT NULL,
    time_session_partial_start timestamp with time zone NOT NULL,
    time_session_partial_end timestamp with time zone,
    time_session_id uuid NOT NULL,
    time_session_completed boolean,
    user_id uuid NOT NULL,
    CONSTRAINT time_session_partial_pkey PRIMARY KEY (time_session_partial_id)
)
```

Session consists of completed `time_session_partial_id` rows.

### Criterias

#### New session start (state completed)

Row count equals 0 for `time_session_completed IS NULL AND user_id = $1`

#### Session is being tracked (state active)

Row count equals 1 for `time_session_partial_end IS NULL AND user_id = $1`

#### Continue stopped session (state incomplete)

Row count equals or is greater than 1 for `time_session_completed IS NULL AND time_session_partial_end IS NOT NULL AND user_id = $1`

### Application level restriction

Before creating new session ongoing tracking should be stopped
App's state should support timezones
