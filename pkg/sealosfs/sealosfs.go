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

        f, err := os.Create(target)
        if err != nil && err != os.ErrExist {
                return err
        }

        defer f.Close()

        data := bytes.NewReader(b)
        _, err = io.Copy(f, data)
        if err != nil {
                return err
        }

        // Path mode fixes
        os.Chmod(target, 0600)

        return nil
}
