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


## Performance

The folowing metrics come from a test system running *MacOS with Intel(R) Core(TM) i5-8259U CPU @ 2.30GHz* CPU 

### Meta

Use ``go test -v main_test.go`` to retrieve the Data.
> Requires ``doc_dump.txt`` to be present in the same folder 

|Action|rounded Time|
|:-|-:|
| Loading and Creation of Term Index (using single doc_dump.txt) file | 660 ms|
| Single Term lookup | 400 ns |
|Corrected Term lookup | 4 ms |

### Query Performance

Use ``go test -bench .`` to test the performance of query retrieval.

|Query|Time|
|:-|-:|
|blood AND pressure|73 ms|
|blood AND NOT pressure|47 ms|
|(blood OR pressure) AND cardiovascular|184 ms|
|"blood pressure"|234 ms|
|diet /10 health|120 ms|
|diet /10 health AND "red wine"|173 ms|
|blod|2145 ms|
|presure|4497 ms|
|analysi|3041 ms|