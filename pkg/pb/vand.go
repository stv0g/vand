// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package pb

import (
	"fmt"
	"io"
	"time"

	"github.com/stv0g/vand/pkg/types"
)

func (s *StateUpdatePoint) Dump(wr io.Writer) {
	f := types.Flatten(s, ".")
	for k, v := range f {
		fmt.Fprintf(wr, "%s = %v\n", k, v) //nolint:errcheck
	}
}

func (ts *Timestamp) Time() time.Time {
	return time.Unix(int64(ts.Seconds), int64(ts.Nanoseconds))
}
