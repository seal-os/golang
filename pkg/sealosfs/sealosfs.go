/*
   Copyright (c) 2018 Open Devices. All rights reserved.
*/

package sealosfs

import (
        "bytes"
        "encoding/json"
        "io"
        "os"
        "sync"

        "github.com/seal-os/sealos-config-api/golang/sealos"
)

var sealOSConflock sync.Mutex

func SaveSealOSConfig(sealos_conf *sealos.APISealOSConfig, target string) error {
        sealOSConflock.Lock()
        defer sealOSConflock.Unlock()

        b, err := json.MarshalIndent(sealos_conf, "", "\t")
        if err != nil {
                return err
        }

        f, err := os.Create(target + ".tmp")
        if err != nil && err != os.ErrExist {
                return err
        }

        data := bytes.NewReader(b)
        _, err = io.Copy(f, data)
        f.Close()
        if err != nil {
                return err
        }

        /* Overwrite final target */
        os.Rename(target + ".tmp", target)

        /* fix permissions */
        os.Chmod(target, 0600)

        return nil
}
