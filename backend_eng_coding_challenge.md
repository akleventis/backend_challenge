# Backend Engineer - Coding Challenge

This coding exercise is an opportunity for you to show us how you break down product requirements into actual code, as well as to demonstrate quality, coding style, experience, and creativity in problem-solving. This task is designed to be relevant to the kind of work an engineer in this role does at Bitly. You should not expect there to be any gotchas.

For the purpose of this challenge, we will be working with CSV and JSON files rather than database tables/streams/APIs but the following will be a representation of data similar to what you would be seeing on a daily basis as an engineer at Bitly.

This unzipped directory contains the data that you will be using for this challenge:
* `encodes.csv` contains information on shortened links or "encodes" to represent existing data infrastructure.
* `decodes.json` contains raw data on bitlink clicks as newline-separated JSON to represent a large data stream.

Use this data set to answer the questions laid out in the problem section of this document.
The data files can be found and downloaded here: https://bit.ly/BitlyBackendCodingChallengeFiles 

## Problem Statement

We would like you to use the provided data files to answer the problem below. Your submission to this challenge should include a runnable program as part of your solution.

**Problem:** Calculate the number of clicks from 2021 for each record in the encode.csv data set.

## Acceptance Criteria

* When run, your program should output the following to the console:
```
A sorted array of JSON objects containing the long URL as the key and the click count as the value. 
The array should be sorted descending by click count.

The form should be as follows: [{"LONG_URL": count}, {"LONG_URL": count}]

Example: [{"https://google.com": 3}, {"https://www.twitter.com" : 2}]
```
* Count clicks only if they occurred in the year 2021
* Solution should include unit tests for each function 
* Solution should be well documented 
* A README which should include:
  - A list of dependencies of your project, as well as how to install them (we may not be experts in your chosen language, framework, and tools)
  - Instructions for running your application/script (you may include a Dockerfile or a Makefile, but this is not a requirement)
  - A description of any design decisions 

## Language and Framework

Feel free to choose any language and framework you are comfortable with. The language that we primarily use at Bitly on the backend is Go, however, you should not feel pressured to use Go if you are not currently comfortable with it. We want you to be able to focus more on your solution than the tools and language.

If you are advanced to the next stage of interviews, the live coding will involve making minor improvements/additions to your coding challenge solution which is why we emphasize using a language that you can navigate comfortably. 

## Notes

* Bitly Glossary:
  * Bitlink: the short link created by the Bitly platform
  * Encode: the act of shortening a long link into a bitlink on the Bitly platform
  * Decode (or click): a metric collected each time a bitlink is accessed and performs the redirect
* The encodes.csv file contains users identified by `user_id` who created a redirect from https://`domain`/`hash` to `long_url`.
