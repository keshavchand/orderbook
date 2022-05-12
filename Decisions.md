# Following are the decisions that were taken when designing the system

* Orders were stored at different price level
* In order matching the traded price will be that of the older order (which was previously present in the book)
* The orders are stored in pro rata basis of order size

## Todo
* Build a reporting system ( report/report.go )
    * Event sourcing
    * Can show the amout held by them
    * Can show the trades in specified time frame ( some sort of timeseries db?? )
    * Calculate the total trading charges for them ( won't be a priority )
