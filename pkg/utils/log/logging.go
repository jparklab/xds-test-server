/*******************************************************************************
*  MIT License
*
*  Copyright (c) 2024 Ji-Young Park(jiyoung.park.dev@gmail.com)
*
*  Permission is hereby granted, free of charge, to any person obtaining a copy
*  of this software and associated documentation files (the "Software"), to deal
*  in the Software without restriction, including without limitation the rights
*  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
*  copies of the Software, and to permit persons to whom the Software is
*  furnished to do so, subject to the following conditions:
*
*      The above copyright notice and this permission notice shall be included in all
*      copies or substantial portions of the Software.
*
*      THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
*      IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
*      FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
*      AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
*      LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
*      OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
*      SOFTWARE.
*******************************************************************************/
package log

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

var (
	defaultLogger *zap.SugaredLogger
)

func init() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}

	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for {
			<-ticker.C
			logger.Sync()
		}
	}()

	defaultLogger = logger.Sugar()
}

func With(args ...interface{}) *zap.SugaredLogger {
	return defaultLogger.With(args...)
}

func Fatal(msg string) {
	defaultLogger.Fatal(msg)
}

func Fatalf(format string, args ...interface{}) {
	defaultLogger.Fatalf(format, args...)
}

func Info(msg string) {
	defaultLogger.Info(msg)
}

func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

func Error(msg string) {
	defaultLogger.Error(msg)
}

func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

func Warn(msg string) {
	defaultLogger.Warn(msg)
}

func Warnf(format string, args ...interface{}) {
	defaultLogger.Warnf(format, args...)
}

func Debug(msg string) {
	defaultLogger.Debug(msg)
}

func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}
