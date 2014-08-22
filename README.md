CronoloGo
=========

An extensible Go library that implements cronolog-like behaviour for os.File descriptors.

Details
-------

Why not log2go?, you may ask. I wanted something simple to handle _just_ opening and closing files for me. I didn't want any log level handling, just "rename the files for me", just as you get when you pipe output through cronolog.

CronoloGo comprises a Rotator object, that will tick away and handle whatever mutations, and the LogFile objects, that describe how you want particular files to be rotated.

The Rotator ticker will take care of the files and issue Close() calls on the old logfile after LogFile.GraceTime (defaults to 0.5s).

You can choose either/both the pointer to File or the callback methods to handle the file references in your code.

The LogFile.Writer property will update the reference in place, so fmt.Fprintf(f, ... ) will work as expected.

With the LogFile.CallBack method, you can wrap the file descriptor in a bufio object, or redirect log.* functions.

Usage Example
-------------
```
func main() {

    var access_log *os.File

    ticker := new(cronologo.Rotator)
    ticker.Start(1 * time.Second)

    ticker.Add(&cronologo.LogFile{
        Writer:     &access_log,
        NamePrefix: "/tmp/log/access.log",
        TimeFormat: "2006-01-02_15:04",
        Symlink:    true,
    })

    ticker.Add(&cronologo.LogFile{
        NamePrefix: "/tmp/log/system.log",
        TimeFormat: "2006-01-02_15:04",
        Symlink:    true,
        GraceTime:  1 * time.Second,
        CallBack:   func(f io.Writer) { log.SetOutput(io.MultiWriter(f, os.Stdout)) },
    })
...

}
```
