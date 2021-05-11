# WPP2-INR

Basic information retrieval system which supports boolean queries.

## Running

Compile the INR-System: `go build .` *(Requires GO >= 1.16)*

Run the INR-System: `./boolret [args]`


### Configuration

Format: `-key value`

The following configuration options are available:

| arg        | desc                                                                                       | default |
|------------|--------------------------------------------------------------------------------------------|---------|
| doc        | Path to the input file (doc_dump.txt)                                                      |         |
| docs       | Path to the input file's                                                                   |         |
| dict       | Path to the dictionary (load if no doc specified - create if doc is specified)             |         |
| k          | size of the k-gram index terms (only relevant if `correction` is `true`)                   | 2       |
| r          | threshold of results before spell correction (only relevant if `correction` is `true`)     | 5       |
| j          | threshold for the jaccard-coefficient  (only relevant if `correction` is `true`)           | 0.2     |
| correction | if spell correction is enabled                                                             | `true`  |

## Queries

* Simple terms
  * `blood`
* Junctions
  * `blood [JUNCTION] pressure`
* Sub-expressions
  * `(blood OR pressure) AND cardiovascular`
* Phrases
  * `"blood pressure"`
* Proximity
  * `diet /[k] health` (eg. `k = 10`)

> **Notice:**
>
> Instead of "\\" for Proximity-Queries we decided to use "/" instead.
>
> eg: ``diet /[k] health`` instead of ``diet \[k] health``

### Junctions

* `AND`
* `AND NOT`
* `OR`

## Benchmarks

Example benchmarks are provided with the _NFCorpus_.
Use ``go test -bench .`` to test the performance of query retrieval.

System Specs:
* OS: Linux
* CPU: Intel(R) Core(TM) i7-7700K CPU @ 4.20GHz

### Meta

Use ``go test -v main_test.go`` to retrieve the Data.

> Requires ``doc_dump.txt`` to be present in the same path.

| action                    | ms |
|---------------------------|---:|
| file read, index creation | 53 |

### Queries

| query                                  |   ns/op |    ms/op |
|----------------------------------------|--------:|---------:|
| blood AND pressure                     |   63533 | 0.063533 |
| blood AND NOT pressure                 |   40837 | 0.040837 |
| (blood OR pressure) AND cardiovascular | 193964  | 0.193964 |
| "blood pressure"                       |  150717 | 0.150717 |
| diet /10 health                        |  102354 | 0.102354 |
| diet /10 health AND "red wine"         |  147999 | 0.147999 |
| blod                                   | 1647286 | 1.647286 |
| presure                                | 3272726 | 3.272726 |
| analysi                                | 2156748 | 2.156748 |
