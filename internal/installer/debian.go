package installer

import (
	"fmt"
	"gokesl/pkg/fs"
	"gokesl/pkg/shell"
	"log"
	"strings"
)

type Debian struct {
	manager            string
	agentpath          string
	keslpath           string
	agentConfigPath    string
	keslConfigPath     string
	keyPath            string
	agentConfigContent []byte
	keslConfigContent  []byte
}

func NewDebian(
	kscIP string,
	basesURI string,
	_keyPath string,
	_agentpath string,
	_keslpath string,
) *Debian {
	return &Debian{
		manager:            dpkg,
		agentpath:          _agentpath,
		keslpath:           _keslpath,
		agentConfigPath:    fs.TempPath(),
		keslConfigPath:     fs.TempPath(),
		keyPath:            _keyPath,
		agentConfigContent: []byte(strings.ReplaceAll(AgentConfig, IpReplacer, kscIP)),
		keslConfigContent:  []byte(strings.ReplaceAll(KeslConfig, KeyReplacer, _keyPath)),
	}
}

func (d *Debian) Install() error {
	if !fs.FileExists("/usr/bin/dpkg") {
		return fmt.Errorf("dpkg not found in PATH")
	}

	exists, err := d.checkPkgExists(klnagent)
	if err != nil {
		return fmt.Errorf("failed to check existing pkg: %w", err)
	}
	if exists {
		if err := d.removePkg(klnagent); err != nil {
			return fmt.Errorf("failed to remove already installed pkg: %w", err)
		}
	}
	if err := d.installPkg(d.agentpath); err != nil {
		return fmt.Errorf("failed to install %s: %w", klnagent, err)
	}
	if err := fs.WriteFile(d.agentConfigPath, d.agentConfigContent); err != nil {
		return fmt.Errorf("failed to create config for %s: %w", klnagent, err)
	}
	if err := d.postinstallAgent(); err != nil {
		return fmt.Errorf("failed to run postinstall agent script: %w", err)
	}
	if err := d.checkAgent(); err != nil {
		return fmt.Errorf("failed to check agent status: %w", err)
	}

	exists, err = d.checkPkgExists(kesl)
	if err != nil {
		return fmt.Errorf("failed to check existing pkg: %w", err)
	}
	if exists {
		if err := d.removePkg(kesl); err != nil {
			return fmt.Errorf("failed to remove already installed pkg: %w", err)
		}
	}
	if err := fs.WriteFile(d.keslConfigPath, d.keslConfigContent); err != nil {
		return fmt.Errorf("failed to create config for %s: %w", kesl, err)
	}
	if err := d.postinstallKesl(); err != nil {
		return fmt.Errorf("failed to run postinstall kesl script: %w", err)
	}

	return nil
}

func (d *Debian) installPkg(pkgpath string) error {
	log.Printf("install pkg: %s\n", pkgpath)
	if _, err := shell.New(d.manager, "--install", pkgpath).Run(); err != nil {
		return err
	}
	return nil
}

func (d *Debian) checkPkgExists(pkgname string) (bool, error) {
	log.Printf("check if %s exists\n", pkgname)
	out, err := shell.New(d.manager, "--list").Run()
	if err != nil {
		return false, err
	}
	if strings.Contains(out, pkgname) {
		return true, nil
	}
	return false, nil
}

func (d *Debian) removePkg(pkgname string) error {
	log.Printf("remove existed pkg %s\n", pkgname)
	if _, err := shell.New(d.manager, "--remove", pkgname).Run(); err != nil {
		return err
	}
	if _, err := shell.New(d.manager, "--purge", pkgname).Run(); err != nil {
		return err
	}
	return nil
}

func (d *Debian) postinstallAgent() error {
	log.Println("run postinstall agent scipt")
	if _, err := shell.New("/opt/kaspersky/klnagent64/lib/bin/setup/postinstall.pl").WithEnv([]string{"KLAUTOANSWERS=" + d.agentConfigPath}).Run(); err != nil {
		return err
	}
	return nil
}

func (d *Debian) checkAgent() error {
	log.Println("run check agent")
	out, code, err := shell.New("/opt/kaspersky/klnagent64/sbin/klnagchk").RunWithReturnCode()
	if code == 1 {
		return nil
	}
	log.Println(out)
	if err != nil {
		return err
	}
	return nil
}

func (d *Debian) postinstallKesl() error {
	log.Println("run postinstall kesl scipt")
	if _, err := shell.New(fmt.Sprintf("/opt/kaspersky/kesl/bin/kesl-setup.pl --autoinstall=%s", d.keslConfigPath)).Run(); err != nil {
		return err
	}
	return nil
}
