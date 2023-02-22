# durt

Disk Utilization Reporting Tool

`durt` produces a simple report of filesystem utilization for provided
directories and files.

With no options it reports on the current directory:

```
% durt
+----------+-------+---------+---------+
| PATH     | FILES |   BYTES |     PCT |
+----------+-------+---------+---------+
| .git     |   105 | 75.74 K |  87.5 % |
| main.go  |     0 |  6.74 K |   7.8 % |
| go.sum   |     0 |  3.11 K |   3.6 % |
| go.mod   |     0 |   578 B | < 1.0 % |
| power.go |     0 |   417 B | < 1.0 % |
+----------+-------+---------+---------+
| TOTALS   |   105 | 86.56 K | 100.0 % |
+----------+-------+---------+---------+
```

It can be used to easily and selectively compare parts of your filesystem that
are using large amounts of disk space.

```
% durt Downloads .Trash Documents
+-----------+-------+----------+---------+
| PATH      | FILES |    BYTES |     PCT |
+-----------+-------+----------+---------+
| Downloads |  1310 |   8.01 G |  68.1 % |
| Documents | 86894 |   3.12 G |  26.5 % |
| .Trash    |   663 | 651.27 M |   5.4 % |
+-----------+-------+----------+---------+
| TOTALS    | 88867 |  11.76 G | 100.0 % |
+-----------+-------+----------+---------+
```
