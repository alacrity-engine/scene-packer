package main

func gameObject(name string) GameObjectPointerMeta {
	return GameObjectPointerMeta{
		Name: name,
	}
}

func resource(typ, id string) ResourcePointerMeta {
	return ResourcePointerMeta{
		ResourceType: typ,
		ResourceID:   id,
	}
}

func component(gmobName, compType string) ComponentPointerMeta {
	return ComponentPointerMeta{
		GmobName: gmobName,
		CompType: compType,
	}
}
