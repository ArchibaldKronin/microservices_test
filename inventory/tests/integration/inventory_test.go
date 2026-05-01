//go:build integration

package integration

import (
	"context"

	inventory_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/inventory/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var _ = Describe("Inventory service", func() {
	var (
		ctx             context.Context
		cancel          context.CancelFunc
		inventoryClient inventory_v1.InventoryServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		conn, err := grpc.NewClient(
			env.App.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешное подключение к gRPC приложению")

		inventoryClient = inventory_v1.NewInventoryServiceClient(conn)
	})

	AfterEach(func() {
		// Чистим коллекцию после теста
		err := env.ClearPartsCollection(ctx)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешную очистку коллекции parts")

		cancel()
	})

	Describe("GetPart", func() {
		var partUUIDsRepo []string
		var partsApi []*inventory_v1.Part

		BeforeEach(func() {
			var err error
			partUUIDsRepo, err = env.InsertTestParts(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку тестовых запчастей в MongoDB")

			partsApi, err = env.GetApiTypeInitialParts(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешную конвертацию исходных данных")
		})

		It("должен успешно возвращать запчасть по ID", func() {
			expPart := partsApi[0]
			id := partUUIDsRepo[0]
			resp, err := inventoryClient.GetPart(ctx, &inventory_v1.GetPartRequest{
				Uuid: id,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetPart()).ToNot(BeNil())
			Expect(resp.GetPart().GetUuid()).To(Equal(id))
			Expect(resp.GetPart().GetManufacturer()).To(Equal(expPart.Manufacturer))
			Expect(resp.GetPart().GetMetadata()).To(Equal(expPart.Metadata))
		})
	})
})
