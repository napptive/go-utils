/**
 * Copyright 2022 Napptive
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package validation

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Param validation tests", func() {

	ginkgo.It("should fail on empty", func() {
		gomega.Expect(CheckNotEmpty("", "name")).To(gomega.HaveOccurred())
	})

	ginkgo.It("should fail on zero", func() {
		gomega.Expect(CheckPositive(0, "name")).To(gomega.HaveOccurred())
	})

	ginkgo.It("should fail on negative values", func() {
		gomega.Expect(CheckPositive(-1, "name")).To(gomega.HaveOccurred())
	})

})
