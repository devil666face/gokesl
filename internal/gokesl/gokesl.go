package gokesl

import (
	"fmt"
	"gokesl/internal/installer"
	"gokesl/pkg/fs"
	"os"
)

const (
	Debian = "debian"
	Redhat = "redhat"
)

type binfile struct {
	data     *[]byte
	filename string
}

type manager interface {
	Install() error
}

type Gokesl struct {
	agentbin  *binfile
	keslbin   *binfile
	os        string
	kscIP     string
	updateURI string
	manager   manager
}

func (b *binfile) Save() error {
	return fs.WriteFile(
		b.filename,
		*b.data,
	)
}

func (b *binfile) Remove() error {
	return os.RemoveAll(b.filename)
}

func New(
	_agentbin, _keslbin *[]byte,
	_os string,
	_kscIP string,
	_updateURI string,
) (*Gokesl, error) {
	_gokesl := &Gokesl{
		agentbin: &binfile{
			data:     _agentbin,
			filename: fs.TempPath(),
		},
		keslbin: &binfile{
			data:     _keslbin,
			filename: fs.TempPath(),
		},
		os:        _os,
		kscIP:     _kscIP,
		updateURI: _updateURI,
	}
	switch _gokesl.os {
	case Debian:
		_gokesl.manager = installer.NewDebian(
			_gokesl.kscIP,
			_gokesl.updateURI,
			_gokesl.agentbin.filename,
			_gokesl.keslbin.filename,
		)
	case Redhat:
		_gokesl.manager = installer.NewRedhat()
	}
	return _gokesl, nil
}

func (g *Gokesl) Install() error {
	var remove func(func() error) = func(f func() error) {
		if err := f(); err != nil {
			fmt.Println(err)
		}
	}
	if err := g.ExtractFiles(); err != nil {
		return fmt.Errorf("failed to extract files: %w", err)
	}
	if err := g.manager.Install(); err != nil {
		return fmt.Errorf("failed to install: %w", err)
	}
	defer remove(g.agentbin.Remove)
	defer remove(g.keslbin.Remove)
	return nil
}

func (g *Gokesl) ExtractFiles() error {
	if err := g.agentbin.Save(); err != nil {
		return err
	}
	if err := g.keslbin.Save(); err != nil {
		return err
	}
	return nil
}
