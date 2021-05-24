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

