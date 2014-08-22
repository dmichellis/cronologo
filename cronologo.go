package cronologo

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Rotator struct {
	loggers []*LogFile
	ticker  *time.Ticker
	running bool
}

func (c *Rotator) Add(l *LogFile) error {
	if err := l.Reopen(); err != nil {
		return err
	}
	if l.GraceTime == 0 {
		l.GraceTime = 1 * time.Millisecond
	}

	loggers_ := append(c.loggers, l)
	c.loggers = loggers_
	return nil
}

func (c *Rotator) Start(d time.Duration) {
	if c.running == true {
		return
	}

	go func() {
		c.ticker = time.NewTicker(d)
		for _ = range c.ticker.C {
			for _, l := range c.loggers {
				go l.Reopen()
			}
		}
	}()
}

func (c *Rotator) Stop() {
	c.ticker.Stop()
}

type LogFile struct {
	oldWriter   *os.File
	Writer      **os.File
	CurrentFile string
	NamePrefix  string
	TimeFormat  string
	Symlink     bool
	CallBack    func(*os.File)
	GraceTime   time.Duration
}

func (c *LogFile) Reopen() error {
	filename := fmt.Sprintf("%s-%s", c.NamePrefix, time.Now().UTC().Format(c.TimeFormat))

	if filename != c.CurrentFile {
		file, err_f := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err_f != nil {
			log.Printf("CronoloGo: Failed to create new logfile '%s': %s - ignoring this change", filename, err_f)
			return err_f
		}

		c.CurrentFile = filename

		if c.Writer != nil {
			(*c.Writer) = file
		}

		if c.Symlink {
			os.Remove(c.NamePrefix)
			err_s := os.Symlink(filename, c.NamePrefix)
			if err_s != nil {
				log.Printf("CronoloGo: Failed to symlink '%s' to '%s': %s", filename, c.NamePrefix, err_s)
			}
		}

		if c.oldWriter != nil {
			oldWriter := c.oldWriter
			go func() {
				time.Sleep(c.GraceTime)
				oldWriter.Close()
			}()
		}
		c.oldWriter = file

		if c.CallBack != nil {
			go c.CallBack(file)
		}

	}
	return nil
}
