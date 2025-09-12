# Memory

## Readonly Files
As a forensic tool, all **write access** to the examined files is **prohibited** by technical measures. All files will be lazy loaded into the main memory by memory mapping them upon the first view.

> To also prevent the writing of configs, caches or plugin outputs to the current filesystem, use the `--readonly` flag or mount the filesystem as readonly. The program will still stay functional.

## Forensic Filesystem
A multi-layered filesystem abstraction is created in-memory upon start. The **base** layer consists of a readonly wrapper around the real filesystem, while to **artifacts** layer holds extracted forensic artifacts of the base file.

## Multicore Operations
All processor heavy operations, like searching or formating, will be done via multicore data handling for faster response times. These operations are optimized for files with one **million** or more lines. 
