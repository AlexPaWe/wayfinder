package cmd
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
	"os"
  "fmt"
  "strings"
  "strconv"
  "runtime"
  "os/signal"

	"github.com/spf13/cobra"
	"github.com/lancs-net/ukbench/log"

	"github.com/lancs-net/ukbench/job"
	"github.com/lancs-net/ukbench/run"
)

type RunConfig struct {
  CpuSets       string
  DryRun        bool
  ScheduleGrace int
}

var (
  runCmd = &cobra.Command{
    Use: "run [OPTIONS...] [FILE]",
    Short: `Run a specific experiment job`,
    Run: doRunCmd,
    Args: cobra.ExactArgs(1),
    DisableFlagsInUseLine: true,
  }
  runConfig = &RunConfig{}
)

func init() {
  runCmd.PersistentFlags().StringVar(
    &runConfig.CpuSets,
    "cpu-sets", 
    fmt.Sprintf("1-%d", runtime.NumCPU()),
    "Specify which CPUs to run experiments on.",
  )
  runCmd.PersistentFlags().BoolVarP(
    &runConfig.DryRun,
    "dry-run",
    "D",
    false,
    "Run without affecting the host or running the jobs.",
  )
  runCmd.PersistentFlags().IntVarP(
    &runConfig.ScheduleGrace,
    "schedule-grace-time",
    "g",
    1,
    "Number of seconds to gracefully wait in the scheduler.",
  )
}

// doRunCmd 
func doRunCmd(cmd *cobra.Command, args []string) {
  // Determine CPU sets
  cpus, err := parseCpuSets(runConfig.CpuSets)
  if err != nil {
    log.Errorf("Could not parse CPU sets: %s", err)
    os.Exit(1)
  }

	j, err := job.NewJob(args[0], &job.RuntimeConfig{
    Cpus:          cpus,
    ScheduleGrace: runConfig.ScheduleGrace,
  })
	if err != nil {
		log.Fatalf("Could not read configuration: %s", err)
		os.Exit(1)
	}

  // Prepare environment
  err = run.PrepareEnvironment(cpus, runConfig.DryRun)
  if err != nil {
    log.Errorf("Could not prepare environment: %s", err)
    cleanup()
    os.Exit(1)
  }

  j.Start(&job.RuntimeConfig{
    Cpus: cpus,
  })

  // We're all done now
  cleanup()
}

func parseCpuSets(cpuSets string) ([]int, error) {
  var cpus []int
  
  if res := strings.Contains(cpuSets, "-"); res {
    c := strings.Split(cpuSets, "-")
    if len(c) > 2 {
      return cpus, fmt.Errorf("Invalid syntax for CPU sets")
    }

    start, err := strconv.Atoi(c[0])
    if err != nil {
      return cpus, fmt.Errorf("Invalid syntax for CPU sets")
    }

    end, err := strconv.Atoi(c[1])
    if err != nil {
      return cpus, fmt.Errorf("Invalid syntax for CPU sets")
    }
    
    for i := start; i < end; i++ {
      cpus = append(cpus, i)
    }
  }

  if strings.Contains(cpuSets, ",") {
    c := strings.Split(cpuSets, ",")

    for i := range c {
      j, err := strconv.Atoi(c[i])
      if err != nil {
        return cpus, fmt.Errorf("Invalid syntax for CPU sets")
      }

      cpus = append(cpus, j)
    }
  }

  return cpus, nil
}

// Create a Ctrl+C trap for reverting machine state
func setupInterruptHandler() {
  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)
  go func(){
    <-c
    cleanup()
    os.Exit(1)
  }()
}

// Preserve the host environment
func cleanup() {
  log.Info("Running clean up...")
  run.RevertEnvironment(runConfig.DryRun)
}
