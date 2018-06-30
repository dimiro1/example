package log

import (
	"github.com/dimiro1/example/toolkit/router"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	logger *logrus.Entry
}

func NewLogger(logger *logrus.Entry) *Logger {
	return &Logger{logger: logger}
}

func (l *Logger) ErrorRendering(err error, handler string) {
	l.logger.WithError(err).WithField("handler", handler).Error("error rendering")
}

func (l *Logger) ErrorOpeningDatabase(err error) {
	l.logger.WithError(err).Fatal("error opening database")
}

func (l *Logger) ErrorCreatingStore(err error, name string) {
	l.logger.WithError(err).WithField("name", name).Fatal("error creating store")
}

func (l *Logger) ErrorCreatingMigrator(err error) {
	l.logger.WithError(err).Fatal("error creating migrator")
}

func (l *Logger) ErrorLoadingTemplates(err error, pattern string) {
	l.logger.WithError(err).WithField("pattern", pattern).Fatal("failed to load templates")
}

func (l *Logger) ErrorInstantiatingModule(err error, name string) {
	l.logger.WithError(err).WithField("module", name).Fatal("failed to create module")
}

func (l *Logger) ErrorCreateApplication(err error) {
	l.logger.WithError(err).Fatal("failed to create application")
}

func (l *Logger) ErrorDatabase(err error) {
	l.logger.WithError(err).Error("error in database")
}

func (l *Logger) ErrorStartingApplication(err error) {
	l.logger.WithError(err).Fatal("failed to start the application")
}

func (l *Logger) StartingApplication() {
	l.logger.Info("starting application")
}

func (l *Logger) RunningMigrations() {
	l.logger.Debug("running migrations")
}

func (l *Logger) FinishedMigrations() {
	l.logger.Debug("finished Running migrations")
}

func (l *Logger) RegisteringRoutes() {
	l.logger.Debug("registering routes")
}

func (l *Logger) FinishedRegisteringRoutes() {
	l.logger.Debug("finished Registering routes")
}

func (l *Logger) Route(route router.Route) {
	l.logger.WithFields(logrus.Fields{
		"method":  route.Method,
		"route":   route.Path,
		"handler": route.HandlerName,
	}).Info()
}

func (l *Logger) ErrorRunningMigrations() {
	l.logger.Info("error running migrations")
}

func (l *Logger) RegisteringModule(name string) {
	l.logger.WithField("module", name).Error("registering module")
}

func (l *Logger) RegisteredModule(name string) {
	l.logger.WithField("module", name).Error("module registered")
}

func (l *Logger) InvalidModule(name string) {
	l.logger.WithField("module", name).Error("could not register module")
}

func (l *Logger) ListeningHTTP() {
	l.logger.Info("listening")
}

func (l *Logger) ErrorListeningHTTP(err error, address string) {
	l.logger.WithError(err).Error("listening HTTP")
}
