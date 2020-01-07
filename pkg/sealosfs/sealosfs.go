/*
   Copyright (c) 2020 Open Devices GmbH. All rights reserved.
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
        b := &bytes.Buffer{}
        encoder := json.NewEncoder(b)
        encoder.SetEscapeHTML(false)
        err := encoder.Encode(sealos_conf)
        if err != nil {
                return err
        }

        out := bytes.Buffer{}
        json.Indent(&out, b.Bytes(), "", "\t")

        f, err := os.OpenFile(target + ".tmp", os.O_RDWR|os.O_CREATE, 0600)
        if err != nil && err != os.ErrExist {
                return err
        }

        data := bytes.NewReader(out.Bytes())
        _, err = io.Copy(f, data)
        if err == nil {
                f.Sync()
        }
        f.Close()
        if err != nil {
                return err
        }

        /* Overwrite final target */
        os.Rename(target + ".tmp", target)

        /* Fix permissions */
        os.Chmod(target, 0600)

        return nil
}
