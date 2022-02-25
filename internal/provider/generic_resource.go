package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/prashantv/tf-test/internal/filestore"
)

type genericResource[T any] struct {
	name        string
	build       func(*schema.ResourceData) T
	set         func(T, *schema.ResourceData) error
	writeStore  func(id string, obj T, create bool) error
	readStore   func(id string) (T, error)
	deleteStore func(id string) error
}

func genID(name string) string {
	t := time.Now()
	return fmt.Sprintf("%v-%v", strings.ToLower(name), t.Format("150405.999999999"))
}

func (r *genericResource[T]) Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := genID(r.name)
	obj := r.build(d)

	if err := r.writeStore(id, obj, true /* create */); err != nil {
		return diag.Errorf("failed to write: %v", err)
	}

	d.SetId(id)
	return nil
}

func (r *genericResource[T]) Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	obj, err := r.readStore(d.Id())
	if err == filestore.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("read %v failed: %v", r.name, err)
	}

	if err := r.set(obj, d); err != nil {
		return diag.Errorf("failed to set: %v", err)
	}

	return nil
}

func (r *genericResource[T]) Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	obj := r.build(d)

	if err := r.writeStore(d.Id(), obj, false /* create */); err != nil {
		return diag.Errorf("failed to write %v: %v", r.name, err)
	}

	return nil
}

func (r *genericResource[T]) Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if err := r.deleteStore(d.Id()); err != nil {
		return diag.Errorf("failed to delete %v: %v", r.name, err)
	}

	return nil
}
