# HealthBridge

Health bridge info

Initial version
    - reads yml config
    - polls health checks at set interval
    - exposes results as prometheus data
    - uses zap for logging

    Test Types:
        - ping test returns true after a 200 response - payload doesnt matter
        - kafdrop test pulls out the lag data from the html and returns that number or -1 on failure

    Next Steps:
        - include internal monitoring tools 
        - run as windows service
        - use cobra/viper to manage params
          - logging - set level

Add pprof for performance profiling 
- go to http://localhost:8080/debug/pprof/ to view performance metrics


Test Coverage:
from - https://blog.golang.org/cover
Generate cover.txt
    go test -coverprofile cover.txt .\...

Generate html report from cover.txt
    go tool cover -html cover.txt -o cover.html


Releases

0.3 Added jeager agent monitor to take a dump of the diagnostic logs in the event of a memory spike

0.5 Updated JaegerAgentMonitor to Memorymonitor
    Changed output to folder containing heap, allocs, version and health info
    
