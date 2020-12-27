package run
// SPDX-License-Identifier: BSD-3-Clause
//
// Authors: Alexander Jung <a.jung@lancs.ac.uk>
//
// Copyright (c) 2020, Lancaster University.  All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
//
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
// 3. Neither the name of the copyright holder nor the names of its
//    contributors may be used to endorse or promote products derived from
//    this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

import (
  "github.com/opencontainers/runc/libcontainer"

  "github.com/lancs-net/ukbench/log"
)

type RuncRunner struct {
  log    *log.Logger
  Config *RunnerConfig
  factory libcontainer.Factory
}

func (r RuncRunner) Init() error {
  // Set the logger
  r.log = r.Config.Log
  
  // Download the image to the cache
  r.log.Infof("Pulling image: %s...", r.Config.Image)
  image, err := PullImage(r.Config.Image, r.Config.CacheDir)
  if err != nil {
    return fmt.Errorf("Could not download image: %s", err)
  }
  
  digest, err := image.Digest()
  if err != nil {
    return fmt.Errorf("Could not process digest: %s", err)
  }
  
  r.log.Debugf("Pulled: %s", digest)

  // Extract the image to the desired location
  r.log.Infof("Extracting image to: %s", r.Config.WorkDir)
  err = UnpackImage(image, r.Config.CacheDir, r.Config.WorkDir, r.Config.AllowOverride)
  if err != nil {
    return fmt.Errorf("Could not extract image: %s", err)
  }

  return nil
}

func (r RuncRunner) Start() error {
  return nil
}

func (r RuncRunner) Wait() (int, error) {
  return 0, nil
}

func (r RuncRunner) Destroy() error {
  return nil
}
