package llm

import "testing"

func TestSplitText(t *testing.T) {
	s := NewMarkDoneSplit(200)
	r := s.Split(`
# bpftune - BPF driven auto-tuning
bpftune aims to provide lightweight, always-on auto-tuning of system behaviour. The key benefit it provides are

by using BPF observability features, we can continuously monitor and adjust system behaviour
because we can observe system behaviour at a fine grain (rather than using coarse system-wide stats), we can tune at a finer grain too (individual socket policies, individual device policies etc)
## Key design principles
- Minimize overhead. Use observability features sparingly; do not trace very high frequency events.
- Be explicit about policy changes providing both a "what" - what change was made - and a "why" - how does it help? syslog logging makes policy actions explicit with explanations
`)
	for i, v := range r {
		t.Logf("%v %v", i, v)
	}

}
