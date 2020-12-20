package job
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

type TaskParam struct {
  Name  string
  Type  string
  Value string
}

// Task is the specific iterated configuration
type Task struct {
  Params   []TaskParam
  Inputs  *[]Input
  Outputs *[]Output
  runs      *Queue
  uuid       string
}

// Init prepare the task 
func (t *Task) Init(runs *[]Run) error {
  // Create a queue of runs for this particular task
  t.runs = NewQueue(len(*runs))

  // Add the runs in-order
  for _, run := range *runs {
    t.runs.Enqueue(run)
  }

  return nil
}

// Cancel the task by removing everything from the queue
func (t *Task) Cancel() {
  log.Warnf("Cancelling task and all subsequent runs")

  // Clear queue of subsequent runs
  t.runs.Clear()
}

func (t *Task) UUID() string {
  if len(t.uuid) == 0 {

    // Calculate the UUID based on a reproducible md5 seed
    md5val := md5.New()
    for _, param := range t.Params {
      io.WriteString(md5val, fmt.Sprintf("%s=%s\n", param.Name, param.Value))
    }

    t.uuid = fmt.Sprintf("%x", md5val.Sum(nil))
  }

  return t.uuid
}

// ActiveTaskRun contains information about a particular task's run.
type ActiveTaskRun struct {
  Task     *Task
  run      *Run
  CoreIds []int // the exact core numbers this task is using
  log      *log.Logger
}
