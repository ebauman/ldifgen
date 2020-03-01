# LDIF Generator

WIP.

Basic tool to generate big-ass LDAP directories.

I was unsatisfied with what I found in the wild to accomplish this task.

Specifically, it seems that everyone has been focused on large user counts but not necessarily complex structures like
deep OUs, many users in many groups, things like that.

So this is my attempt to address those needs.

## Usage

Outputs to `stdout`. Pipe, redirect, do whatever you need to. 

```text
$ ./ldifgen generate --help
NAME:
   ldifgen generate - generate ldif file

USAGE:
   ldifgen generate [command options] [arguments...]

OPTIONS:
   --users value               number of users to generate (default: 10)
   --ous value                 number of organizational units to generate (default: 2)
   --ou-depth value            depth of generated OUs. specify n>1 to create 'chains' of OUs (default: 1)
   --groups value              number of groups to generate (default: 2)
   --domain value              domain used to generate DC components, e.g. dc=domain,dc=example,dc=org (default: "domain.example.org")
   --user-classes value        comma-separated list of classes for user objects (default: "top,person,organizationalPerson,inetOrgPerson")
   --ou-classes value          comma-separated list of classes for organizational unit objects (default: "top,organizationalUnit")
   --group-classes value       comma-separated list of classes for group objects (default: "top,groupOfNames")
   --user-change-type value    LDIF changetype for users (default: "add")
   --group-change-type value   LDIF changetype for groups (default: "add")
   --ou-change-type value      LDIF changetype for OUs (default: "add")
   --buzzword-dataset value    path to an alternative list of buzzwords, used in group generation. provide list of words, separated by newlines
   --department-dataset value  path to an alternative list of department names, used in OU generation. provide list of words, separated by newlines
   --first-name-dataset value  path to an alternative list of first names, used in user generation. provide list of words, separated by newlines
   --last-name-dataset value   path to an alternative list of last names, used in user generation. provide list of words, separated by newlines
   --groups-dataset value      path to an alternative list of group names, used in group generation. provide list of words, separated by newlines
   --help, -h                  show help (default: false)

```

## Features

- Random name generation for all objects. Pre-seeded with templates as seen in `static/` directory, but you can supply
your own via the `--xxx-dataset` arguments. 
- Prevents name collisions via retry logic. Very basic, just cycles until hitting 5 collisions, then bails. 
- Generates simple OU trees, which if seen graphically are really just a single linked list currently. Room for improvement
there. 
- Generates groups and adds even weighting of users to each group. Future goal: user-definable weighting.
- Customizable change types per object
- Customizable object classes for OUs, groups, users

## Performance

It can build a 100k-user, 1000-group, 20-OU (at a depth of 20 sub-OUs) in 10 seconds.

```text
$ time ./ldifgen generate --users 100000 --groups 1000 --ous 20 --ou-depth 20 > out.ldif
INFO[0000] OU name collision: Finance. Generating new OU 
INFO[0000] OU name collision: Arts and Sciences. Generating new OU 
INFO[0000] OU name collision: Center for Computational Science and Engineering. Generating new OU 
INFO[0000] OU name collision: Anthropology. Generating new OU 
INFO[0000] group name collision. generating new group   

real    0m10.994s
user    0m4.061s
sys     0m7.278s

```

## To-Do
- Allow for creation of custom object types? 
- Allow for weighting of users across groups (currently just weights evenly)
- Allow for broader OU trees, currently they're the equivalent of a linked list
- Custom templates?