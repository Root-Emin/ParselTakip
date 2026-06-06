-- +goose Up
-- +goose StatementBegin

-- Sistem Rolleri
INSERT INTO system_role_definitions (code, name, description, permissions) VALUES
('municipality_admin', 'Belediye Yöneticisi', 'Kentsel dönüşüm başlatır, bina oluşturur, müteahhit atar, süreçleri denetler, onay verir, raporları görüntüler',
 '["project.create","project.approve","building.create","contractor.assign","report.view","workflow.advance"]'),
('municipality_staff', 'Belediye Personeli', 'Evrak kontrolü yapar, eksik evrak takibi yapar, vatandaşlarla iletişim kurar',
 '["document.review","document.track","citizen.contact","workflow.update"]'),
('contractor', 'Müteahhit Firma', 'Atandığı projeleri görür, evrak yükler, süreç günceller, hak sahipleri ile iletişim kurar',
 '["project.view_assigned","document.upload","progress.update","owner.contact"]'),
('property_owner', 'Hak Sahibi (Vatandaş)', 'Evrak yükler, süreç durumunu görür, bildirim alır, onay verir',
 '["document.upload","process.view","notification.receive","consent.give","objection.file"]'),
('system_admin', 'Sistem Yöneticisi', 'Belediye yönetimi, kullanıcı yönetimi, yetki yönetimi, sistem ayarları',
 '["organization.manage","user.manage","role.manage","system.configure"]')
ON CONFLICT (code) DO NOTHING;

-- Belge Türleri (Gerekli + Geçersiz olanlar)
INSERT INTO document_types (code, name, description, category, is_mandatory, requires_notary, invalid_reason) VALUES
('risk_assessment_report', 'Riskli Yapı Tespit Raporu', 'Yetkili kurum/kuruluş tarafından hazırlanan risk raporu', 'technical', TRUE, FALSE, NULL),
('title_deed', 'Tapu Belgesi', 'Kat irtifakı veya kat mülkiyeti tapusu', 'legal', TRUE, FALSE, NULL),
('identity_document', 'Kimlik Belgesi', 'Maliklerin kimlik fotokopileri ve ikamet bilgileri', 'identity', TRUE, FALSE, NULL),
('construction_contract', 'Kat Karşılığı İnşaat Sözleşmesi', 'Noter onaylı müteahhit sözleşmesi', 'contract', TRUE, TRUE, 'Noter onayı olmayan sözleşmeler geçersizdir'),
('building_permit', 'Yapı Ruhsatı', 'Belediye yapı ruhsat belgesi', 'permit', TRUE, FALSE, 'Belediye onayı olmayan projeler geçersizdir'),
('zoning_compliance', 'İmar Planı Uygunluk Belgesi', 'İmar planına uygunluk belgesi', 'permit', TRUE, FALSE, NULL),
('occupancy_permit', 'İskan Belgesi', 'Yapı kullanma izin belgesi', 'permit', TRUE, FALSE, NULL),
('tax_exemption', 'Vergi ve Harç Muafiyet Belgesi', 'Tapu harcı, damga vergisi, noter harcı muafiyet belgeleri', 'tax', FALSE, FALSE, NULL),
('rent_assistance_application', 'Kira Yardımı Başvuru Dilekçesi', 'Çevre ve Şehircilik İl Müdürlüğü başvuru belgesi', 'financial', FALSE, FALSE, NULL),
('lease_contract', 'Kira Kontratı', 'Kira yardımı için kira sözleşmesi', 'financial', FALSE, FALSE, NULL),
('informal_agreement', 'Gayriresmi Anlaşma', 'Yazılı olmayan mutabakatlar', 'other', FALSE, FALSE, 'Gayriresmi anlaşmalar geçersizdir ve kullanılamaz'),
('unsigned_contract', 'Eksik İmzalı Sözleşme', 'Noter onaysız kat karşılığı sözleşme', 'contract', FALSE, FALSE, 'Noter onayı olmayan sözleşmeler geçersizdir')
ON CONFLICT (code) DO NOTHING;

-- İş Akışı Adımları (9 adım)
INSERT INTO workflow_step_definitions (step_order, code, name, description, responsible_role, sla_days) VALUES
(1, 'risk_assessment', 'Riskli Yapı Tespiti', 'Yetkili kurum risk raporu hazırlar', 'municipality_admin', 30),
(2, 'land_registry_notification', 'Tapu Müdürlüğüne Bildirim', 'Riskli yapı tapu müdürlüğüne bildirilir', 'municipality_staff', 7),
(3, 'owner_notification', 'Maliklere Tebligat', '15 gün içinde itiraz hakkı tanınır', 'municipality_staff', 15),
(4, 'majority_decision', '2/3 Çoğunluk Kararı', 'Müteahhit seçimi ve sözleşme kararı', 'property_owner', 30),
(5, 'contract_signing', 'Kat Karşılığı Sözleşme', 'Noter onaylı sözleşme imzalanır', 'contractor', 15),
(6, 'permit_process', 'Belediye Ruhsat Süreci', 'İmar planı uygunluk ve yapı ruhsatı alınır', 'municipality_admin', 45),
(7, 'demolition', 'Yıkım İşlemi', 'Eski bina yıkılır, inşaata başlanır', 'contractor', 30),
(8, 'rent_assistance', 'Kira Yardımı Süreci', 'Malikler başvurur, devlet ödeme yapar', 'property_owner', 60),
(9, 'construction_delivery', 'İnşaat ve Teslim', 'Proje tamamlanır, iskan alınır, yeni tapular verilir', 'contractor', 365)
ON CONFLICT (code) DO NOTHING;

-- Adım bazlı belge gereksinimleri
INSERT INTO workflow_document_requirements (workflow_step_id, document_type_id, is_mandatory, responsible_role, notes)
SELECT ws.id, dt.id, TRUE, 'municipality_admin', 'Risk raporu olmadan süreç başlamaz'
FROM workflow_step_definitions ws, document_types dt
WHERE ws.code = 'risk_assessment' AND dt.code = 'risk_assessment_report'
ON CONFLICT DO NOTHING;

INSERT INTO workflow_document_requirements (workflow_step_id, document_type_id, is_mandatory, responsible_role, notes)
SELECT ws.id, dt.id, TRUE, 'municipality_staff', 'Tapu müdürlüğüne bildirim için gerekli'
FROM workflow_step_definitions ws, document_types dt
WHERE ws.code = 'land_registry_notification' AND dt.code IN ('risk_assessment_report', 'title_deed')
ON CONFLICT DO NOTHING;

INSERT INTO workflow_document_requirements (workflow_step_id, document_type_id, is_mandatory, responsible_role, notes)
SELECT ws.id, dt.id, TRUE, 'property_owner', 'Malik kimlik belgeleri'
FROM workflow_step_definitions ws, document_types dt
WHERE ws.code = 'owner_notification' AND dt.code = 'identity_document'
ON CONFLICT DO NOTHING;

INSERT INTO workflow_document_requirements (workflow_step_id, document_type_id, is_mandatory, responsible_role, notes)
SELECT ws.id, dt.id, TRUE, 'contractor', 'Noter onaylı olmalı'
FROM workflow_step_definitions ws, document_types dt
WHERE ws.code = 'contract_signing' AND dt.code = 'construction_contract'
ON CONFLICT DO NOTHING;

INSERT INTO workflow_document_requirements (workflow_step_id, document_type_id, is_mandatory, responsible_role, notes)
SELECT ws.id, dt.id, TRUE, 'municipality_admin', 'Ruhsat olmadan inşaat yapılamaz'
FROM workflow_step_definitions ws, document_types dt
WHERE ws.code = 'permit_process' AND dt.code IN ('building_permit', 'zoning_compliance')
ON CONFLICT DO NOTHING;

INSERT INTO workflow_document_requirements (workflow_step_id, document_type_id, is_mandatory, responsible_role, notes)
SELECT ws.id, dt.id, FALSE, 'property_owner', 'Kira yardımı başvurusu'
FROM workflow_step_definitions ws, document_types dt
WHERE ws.code = 'rent_assistance' AND dt.code IN ('rent_assistance_application', 'lease_contract')
ON CONFLICT DO NOTHING;

INSERT INTO workflow_document_requirements (workflow_step_id, document_type_id, is_mandatory, responsible_role, notes)
SELECT ws.id, dt.id, TRUE, 'municipality_admin', 'İnşaat sonunda iskan alınır'
FROM workflow_step_definitions ws, document_types dt
WHERE ws.code = 'construction_delivery' AND dt.code IN ('occupancy_permit', 'title_deed')
ON CONFLICT DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM workflow_document_requirements;
DELETE FROM workflow_step_definitions;
DELETE FROM document_types;
DELETE FROM system_role_definitions;
-- +goose StatementEnd
