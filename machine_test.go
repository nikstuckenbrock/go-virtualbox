package virtualbox

import (
	"context"
	"errors"
	"testing"

	"github.com/go-test/deep"
	"github.com/golang/mock/gomock"
)

var (
	testUbuntuMachine = &Machine{
		Name:     "Ubuntu",
		Firmware: "BIOS",
		UUID:     "37f5d336-bf07-48dd-947c-37e6a56420a7",
		State:    Saved,
		CPUs:     1,
		Memory:   1024, VRAM: 8, CfgFile: "/Users/fix/VirtualBox VMs/go-virtualbox/go-virtualbox.vbox",
		BaseFolder: "/Users/fix/VirtualBox VMs/go-virtualbox", OSType: "", Flag: 0, BootOrder: []string{},
		NICs: []NIC{
			{Network: "nat", Hardware: "82540EM", HostInterface: "", MacAddr: "080027EE1DF7"},
		},
	}
)

func TestMachine(t *testing.T) {
	testCases := map[string]struct {
		in   string
		want *Machine
		err  error
	}{
		"by name": {
			in:   "Ubuntu",
			want: testUbuntuMachine,
			err:  nil,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			m := newTestManager()

			got, err := m.Machine(context.Background(), tc.in)
			if diff := deep.Equal(got, tc.want); !errors.Is(err, tc.err) || diff != nil {
				t.Errorf("Machine(%s) = %+v, %v; want %v, %v; diff = %v",
					tc.in, got, err, tc.want, tc.err, diff)
			}
		})
	}
}

func TestLegacy_Machine(t *testing.T) {
	Setup(t)

	if ManageMock != nil {
		listVmsOut := ReadTestData("vboxmanage-list-vms-1.out")
		vmInfoOut := ReadTestData("vboxmanage-showvminfo-1.out")
		gomock.InOrder(
			ManageMock.EXPECT().runOut("list", "vms").Return(listVmsOut, nil).Times(1),
			ManageMock.EXPECT().runOutErr("showvminfo", "Ubuntu", "--machinereadable").Return(vmInfoOut, "", nil).Times(1),
			ManageMock.EXPECT().runOutErr("showvminfo", "go-virtualbox", "--machinereadable").Return(vmInfoOut, "", nil).Times(1),
		)
	}
	ms, err := ListMachines()
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range ms {
		t.Logf("%+v", m)
	}

	Teardown()
}
