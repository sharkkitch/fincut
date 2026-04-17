// Package rotate detects log file rotation by monitoring inode changes
// or file truncation, and notifies callers so they can reopen the file.
//
// Basic usage:
//
//	rot, err := rotate.NewRotator(rotate.Options{
//		Path:     "/var/log/app.log",
//		Output:   os.Stdout,
//		Interval: 2 * time.Second,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//	rot.Run(ctx)
package rotate
