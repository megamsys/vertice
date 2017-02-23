package marketplaces

const (
  //actions available for rawimage.
  RAW_ISO_CREATE = "rawimage.iso.create"

  // actions available for marketplaces
  MARKETPLACE_ISO_CUSTOMIZE = "marketplaces.iso.customize"
  MARKETPLACE_ISO_FINISHED = "marketplaces.iso.finished"  // made backup of vm (save changes to Datablock)
  MARKETPLACE_IMAGE_ADD =   "marketplaces.image.add"
)

type MarketplaceInterface interface {
   Process(action string) error
   String() string
}

//process trigger based on acction
func (m *Marketplaces) Process(action string) error {
  switch action {
  case MARKETPLACE_ISO_CUSTOMIZE:
   return m.rawImageCustomize()
  default:
    return nil
  }
  return nil
}

//process trigger based on acction
func (s *RawImages) Process(action string) error {
  switch action {
  case RAW_ISO_CREATE:
   return s.create()
  default:
  }
  return nil
}
