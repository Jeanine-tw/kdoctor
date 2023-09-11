
# environment
- Kubenetes: `v1.27.1`
- container runtime: `containerd 1.6.12`
- OS: `CentOS Linux 8`
- kernel: `4.18.0-348.7.1.el8_5.x86_64`

| Node     | Role          | CPU  | Memory |
| -------- | ------------- | ---- | ------ |
| master1  | control-plane | 4C   | 8Gi    |
| master2  | control-plane | 4C   | 8Gi    |
| master3  | control-plane | 4C   | 8Gi    |
| worker4  |               | 3C   | 8Gi    |
| worker5  |               | 3C   | 8Gi    |
| worker6  |               | 3C   | 8Gi    |
| worker7  |               | 3C   | 8Gi    |
| worker8  |               | 3C   | 8Gi    |
| worker9  |               | 3C   | 8Gi    |
| worker10 |               | 3C   | 8Gi    |

# Nethttp

In a pod with a CPU of 1

The test server is a server that sleeps for one second and then returns

The client uses 50 concurrency

## Http1.1

| client  | time | requests | qps   | Memory |
|---------|------|----------|-------|--------|
| kdoctor | 0.5m | 63157    | 2105  | 319Mb  |
| ab      | 0.5m | 1402     | 46.62 | 2Mb    |
| wrk     | 0.5m | 1472     | 48.99 | 7Mb    |
| hey     | 0.5m | 1495     | 49.84 | 33.5Mb |


| client  | time | requests  | qps   | Memory |
|---------|------|-----------|-------|--------|
| kdoctor | 1m   | 115825    | 1930  | 502Mb  |
| ab      | 1m   | 2902      | 48.27 | 2Mb    |
| wrk     | 1m   | 2957      | 49.21 | 7Mb    |
| hey     | 1m   | 2991      | 49.85 | 53.5Mb |

| client  | time | requests | qps   | Memory |
|---------|------|----------|-------|--------|
| kdoctor | 5m   | 598500   | 1995  | 630Mb  |
| ab      | 5m   | 14902    | 49.56 | 2.5Mb  |
| wrk     | 5m   | 14950    | 49.82 | 7.5Mb  |
| hey     | 5m   | 14958    | 49.86 | 80.5Mb |


Data reporting when kdoctor turns metric on and off

| client  | time | requests | qps  | Memory | enable metric |
|---------|------|----------|------|--------|---------------|
| kdoctor | 0.5m | 60900    | 2030 | 431Mb  | true          |
| kdoctor | 0.5m | 63157    | 2105 | 319Mb  | false         |
| kdoctor | 1m   | 99750    | 1995 | 530Mb  | true          |
| kdoctor | 1m   | 115825   | 1930 | 502Mb  | false         |
| kdoctor | 5m   | 561000   | 1870 | 827Mb  | true          |
| kdoctor | 5m   | 598500   | 1995 | 630Mb  | false         |
| kdoctor | 10m  | 1158000  | 1930 | 780Mb  | true          |
| kdoctor | 10m  | 1122000  | 1870 | 730Mb  | false         |
| kdoctor | 15m  | 1788000  | 1986 | 623Mb  | true          |
| kdoctor | 15m  | 1800000  | 2000 | 673Mb  | false         |
| kdoctor | 20m  | 2253600  | 1878 | 752Mb  | true          |
| kdoctor | 20m  | 2256000  | 1880 | 739Mb  | false         |
| kdoctor | 30m  | 3332000  | 1851 | 742Mb  | true          |
| kdoctor | 30m  | 3112000  | 1729 | 698Mb  | false         |


## Http2

| client  | time | requests | qps     | Memory |
|---------|------|----------|---------|--------|
| kdoctor | 0.5m | 238787   | 7959.57 | 350Mb  |
| hey     | 0.5m | 7213     | 240.44  | 110Mb  |

| client  | time | requests | qps      | Memory |
|---------|------|----------|----------|--------|
| kdoctor | 1m   | 481070   | 8017.83  | 370Mb  |
| hey     | 1m   | 14665    | 244.42   | 120Mb  |

| client  | time | requests | qps      | Memory |
|---------|------|----------|----------|--------|
| kdoctor | 5m   | 2419874  | 8066.25  | 390Mb  |
| hey     | 5m   | 74776    | 249.25   | 130Mb  |


# Netdns

In a pod with a CPU of 1

| client  | time | requests | qps      | Memory |
|---------|------|----------|----------|--------|
| kdoctor | 1m   | 1921378  | 32022.97 | 42Mb   |
| dnsperf | 1m   | 1728086  | 28800.40 | 8Mb    |

| client  | time | requests | qps      | Memory |
|---------|------|----------|----------|--------|
| kdoctor | 5m   | 9599260  | 31997.53 | 42Mb   |
| dnsperf | 5m   | 8811137  | 29370.34 | 8Mb    |

| client  | time | requests  | qps      | Memory |
|---------|------|-----------|----------|--------|
| kdoctor | 10m  | 19166303  | 31943.84 | 47Mb   |
| dnsperf | 10m  | 17260779  | 28767.66 | 8Mb    |

Data reporting when kdoctor turns metric on and off

| client  | time | requests | qps      | Memory | enable metric |
|---------|------|----------|----------|--------|---------------|
| kdoctor | 1m   | 1910392  | 31839.87 | 71Mb   | true          |
| kdoctor | 1m   | 1921378  | 32022.97 | 42Mb   | false         |
| kdoctor | 5m   | 9651682  | 32172.27 | 194Mb  | true          |
| kdoctor | 5m   | 9599260  | 31997.53 | 42Mb   | false         |
| kdoctor | 10m  | 19446976 | 32411.63 | 366Mb  | true          |
| kdoctor | 10m  | 19166303 | 31943.84 | 47Mb   | false         |
| kdoctor | 15m  | 29224060 | 32471.18 | 513Mb  | true          |
| kdoctor | 15m  | 28836875 | 32040.97 | 46Mb   | false         |
| kdoctor | 20m  | 38431238 | 32026.03 | 584Mb  | true          |
| kdoctor | 20m  | 38485112 | 32070.93 | 42Mb   | false         |